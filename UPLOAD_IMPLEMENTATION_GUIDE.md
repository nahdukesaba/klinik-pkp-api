# Upload Implementation Guide

## Summary of Changes

I've refactored the upload system to support **multiple image uploads (1-4 images)** directly within each module's create/update endpoints.

### New Folder Structure
```
images/
  └── {category}/          # e.g., sosialisasi, rusun, kumuh, bsps
      └── YYYY/
          └── MM/
              └── DD/
                  └── {record_id}/     # NEW: Groups images by record
                      └── {category}_{uuid}.{ext}
```

**Benefits:**
- Easy to delete all images when deleting a record
- Organizes images by date for archival/cleanup
- Prevents mixing images from different records

---

## ✅ Completed: Sosialisasi Module

The **sosialisasi** module is fully implemented with:
- ✅ Multiple image support (1-4 images)
- ✅ Image upload during create
- ✅ Image replacement during update
- ✅ Auto-deletion of images on record delete
- ✅ Multipart form handling

---

## 🔧 To Implement: Rusun, Kumuh, BSPS Modules

Follow these steps for each remaining module:

### Step 1: Update Service Constructor

**Before:**
```go
func NewService(db *gorm.DB) *Service {
    return &Service{db: db, validator: &utils.Validator{}}
}
```

**After:**
```go
import "klinik-pkp-api/internal/api/uploads"

type Service struct {
    db            *gorm.DB
    validator     *utils.Validator
    uploadService *uploads.Service  // ADD THIS
}

func NewService(db *gorm.DB, uploadService *uploads.Service) *Service {
    return &Service{
        db:            db,
        validator:     &utils.Validator{},
        uploadService: uploadService,  // ADD THIS
    }
}
```

### Step 2: Update Payload Struct

**For modules with numeric IDs (rusun, kumuh, bsps):**
```go
import "mime/multipart"

type RusunPayload struct {
    VillageID   string                  `form:"village_id"`
    DistrictID  string                  `form:"district_id"`
    RegionID    string                  `form:"region_id"`
    Name        string                  `form:"name"`
    Address     string                  `form:"address"`
    // ... other fields ...
    Images      []*multipart.FileHeader `form:"images"`  // ADD THIS
}
```

### Step 3: Update Create Method

Add image upload logic after creating the record:

```go
func (s *Service) AddRusun(payload *RusunPayload) (*Rusun, error) {
    // ... existing validation code ...

    // CREATE RECORD FIRST
    rusun := Rusun{
        VillageID: payload.VillageID,
        // ... other fields ...
    }

    if err := s.db.Create(&rusun).Error; err != nil {
        return nil, err
    }

    // HANDLE IMAGE UPLOADS if provided
    if len(payload.Images) > 0 {
        imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
            Files:    payload.Images,
            Category: "rusun",  // CHANGE TO YOUR MODULE NAME
            RecordID: fmt.Sprintf("%d", rusun.ID),  // Convert uint to string
        })

        if err != nil {
            // Rollback: delete the created record
            s.db.Delete(&rusun)
            return nil, fmt.Errorf("failed to upload images: %w", err)
        }

        // Extract URLs from responses
        var imageURLs []string
        for _, img := range imageResponses {
            imageURLs = append(imageURLs, img.URL)
        }

        // Update the record with image URLs
        rusun.ImageURLs = imageURLs
        if err := s.db.Save(&rusun).Error; err != nil {
            return nil, err
        }
    }

    return &rusun, nil
}
```

### Step 4: Add Update Method

```go
func (s *Service) UpdateRusun(id uint, payload *RusunPayload) (*Rusun, error) {
    // FIND EXISTING RECORD
    var rusun Rusun
    if err := s.db.First(&rusun, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.New("Rusun not found")
        }
        return nil, err
    }

    // UPDATE FIELDS (only if provided)
    if payload.VillageID != "" {
        rusun.VillageID = payload.VillageID
    }
    // ... update other fields ...

    // HANDLE IMAGE UPDATES
    if len(payload.Images) > 0 {
        // Delete old images
        if err := s.uploadService.DeleteImages("rusun", fmt.Sprintf("%d", id)); err != nil {
            fmt.Printf("Warning: failed to delete old images: %v\n", err)
        }

        // Upload new images
        imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
            Files:    payload.Images,
            Category: "rusun",
            RecordID: fmt.Sprintf("%d", id),
        })

        if err != nil {
            return nil, fmt.Errorf("failed to upload new images: %w", err)
        }

        // Extract URLs
        var imageURLs []string
        for _, img := range imageResponses {
            imageURLs = append(imageURLs, img.URL)
        }
        rusun.ImageURLs = imageURLs
    }

    // SAVE UPDATES
    if err := s.db.Save(&rusun).Error; err != nil {
        return nil, err
    }

    return &rusun, nil
}
```

### Step 5: Add Delete Method

```go
func (s *Service) DeleteRusun(id uint) error {
    var rusun Rusun
    if err := s.db.First(&rusun, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return errors.New("Rusun not found")
        }
        return err
    }

    // DELETE ASSOCIATED IMAGES
    if err := s.uploadService.DeleteImages("rusun", fmt.Sprintf("%d", id)); err != nil {
        fmt.Printf("Warning: failed to delete images: %v\n", err)
    }

    // DELETE RECORD
    if err := s.db.Delete(&rusun).Error; err != nil {
        return err
    }

    return nil
}
```

### Step 6: Update Handler

```go
func (h *Handler) AddRusunHandler(ctx *fiber.Ctx) error {
    // Parse multipart form
    form, err := ctx.MultipartForm()
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
    }

    // Create payload
    payload := &RusunPayload{
        VillageID:  ctx.FormValue("village_id"),
        DistrictID: ctx.FormValue("district_id"),
        RegionID:   ctx.FormValue("region_id"),
        Name:       ctx.FormValue("name"),
        // ... other fields using ctx.FormValue() ...
        Images:     form.File["images"],
    }

    rusun, err := h.service.AddRusun(payload)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
    }

    return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse("Rusun created successfully", rusun))
}

func (h *Handler) UpdateRusunHandler(ctx *fiber.Ctx) error {
    id, err := ctx.ParamsInt("id")
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
    }

    form, err := ctx.MultipartForm()
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
    }

    payload := &RusunPayload{
        VillageID: ctx.FormValue("village_id"),
        // ... other fields ...
        Images:    form.File["images"],
    }

    rusun, err := h.service.UpdateRusun(uint(id), payload)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
    }

    return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Rusun updated successfully", rusun))
}

func (h *Handler) DeleteRusunHandler(ctx *fiber.Ctx) error {
    id, err := ctx.ParamsInt("id")
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
    }

    err = h.service.DeleteRusun(uint(id))
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
    }

    return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Rusun deleted successfully", nil))
}
```

### Step 7: Update main.go

For **each module**, update the service initialization:

```go
// INIT SERVICE FOR EACH MODULE
uploadsService := _uploads.NewService("./images")
// ...
rusunService := rusun.NewService(db, uploadsService)  // ADD uploadsService
kumuhService := kumuh.NewService(db, uploadsService)  // ADD uploadsService
bspsService := bsps.NewService(db, uploadsService)    // ADD uploadsService
```

---

## 📝 Frontend Usage

### Creating a Record with Images

```javascript
const formData = new FormData();
formData.append('village_id', '123');
formData.append('district_id', '456');
formData.append('title', 'Event Title');
formData.append('location', 'Event Location');
formData.append('description', 'Description here');
formData.append('scheduled_at', '2026-01-15T10:00:00Z');

// Add multiple images
const images = document.querySelector('input[type="file"]').files;
for (let i = 0; i < images.length; i++) {
    formData.append('images', images[i]);
}

const response = await fetch('/api/v1/sosialisasi', {
    method: 'POST',
    body: formData,
    headers: {
        'Authorization': 'Bearer YOUR_TOKEN'
    }
});

const data = await response.json();
console.log(data.data.image_urls); // Array of image URLs
```

### Updating with New Images

```javascript
const formData = new FormData();
formData.append('title', 'Updated Title');

// Upload new images (replaces old ones)
const newImages = document.querySelector('input[type="file"]').files;
for (let i = 0; i < newImages.length; i++) {
    formData.append('images', newImages[i]);
}

await fetch('/api/v1/sosialisasi/abc-123', {
    method: 'PUT',
    body: formData
});
```

---

## 🗄️ Database Migration

You need to add the `image_urls` column to each table:

```sql
-- Sosialisasi
ALTER TABLE sosialisasi 
  DROP COLUMN image_url,
  ADD COLUMN image_urls JSON;

-- Rusun
ALTER TABLE rusun 
  ADD COLUMN image_urls JSON;

-- Kawasan Kumuh
ALTER TABLE kawasan_kumuh 
  ADD COLUMN image_urls JSON;

-- BSPS
ALTER TABLE penerimaan_bsps 
  ADD COLUMN image_urls JSON;
```

Or create a migration file in your `migrations/` folder.

---

## 🚀 Next Steps

1. ✅ Sosialisasi is done
2. ⬜ Implement Rusun (follow steps above)
3. ⬜ Implement Kumuh (follow steps above)
4. ⬜ Implement BSPS (follow steps above)
5. ⬜ Run database migrations
6. ⬜ Test with Postman/frontend

---

## 📌 Key Points

- **Max 4 images** per record
- Images stored as JSON array in database
- Folder structure: `images/{category}/YYYY/MM/DD/{record_id}/filename.ext`
- Old images automatically deleted on update/delete
- Supports JPEG, PNG, GIF, WebP (configurable in `utils/image.go`)
- Max file size: 10MB (configurable in `utils/image.go`)

---

Need help implementing any specific module? Let me know!
