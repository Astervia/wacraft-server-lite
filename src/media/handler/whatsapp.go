package media_handler

import (
	"errors"
	"strconv"

	cmn_model "github.com/Astervia/wacraft-core/src/common/model"
	common_service "github.com/Astervia/wacraft-core/src/common/service"
	"github.com/Astervia/wacraft-server/src/integration/whatsapp"
	common_model "github.com/Rfluid/whatsapp-cloud-api/src/common/model"
	media_model "github.com/Rfluid/whatsapp-cloud-api/src/media/model"
	media_service "github.com/Rfluid/whatsapp-cloud-api/src/media/service"
	"github.com/gofiber/fiber/v2"
)

// GetWhatsAppMediaURL retrieves a temporary download URL for a WhatsApp media item.
//	@Summary		Gets URL for WhatsApp media
//	@Description	Uses the WhatsApp API to get a temporary URL to download the media. The URL expires in 5 minutes.
//	@Tags			Media
//	@Accept			json
//	@Produce		json
//	@Param			mediaId	path		string							true	"Media ID"
//	@Success		200		{object}	media_model.MediaInfo			"Media information with download URL"
//	@Failure		400		{object}	common_model.DescriptiveError	"Missing or invalid media ID"
//	@Failure		500		{object}	common_model.DescriptiveError	"Failed to retrieve media URL"
//	@Security		ApiKeyAuth
//	@Router			/media/whatsapp/{mediaId} [get]
func GetWhatsAppMediaURL(ctx *fiber.Ctx) error {
	mediaId := ctx.Params("mediaId")
	if mediaId == "" {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			cmn_model.NewApiError("media ID is required", errors.New("no media ID provided"), "handler").Send(),
		)
	}

	mediaInfo, err := media_service.RetrieveURL(whatsapp.WabaApi, mediaId, media_model.RetrieveInfo{})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			cmn_model.NewApiError("failed to retrieve media URL", err, "handler").Send(),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(mediaInfo)
}

// DownloadWhatsAppMedia downloads a media file directly from WhatsApp using its media ID.
//	@Summary		Downloads WhatsApp media
//	@Description	Downloads media using the URL retrieved via the WhatsApp API.
//	@Tags			Media
//	@Accept			json
//	@Produce		application/octet-stream
//	@Param			mediaId	path		string							true	"Media ID"
//	@Success		200		{file}		binary							"Downloaded media file"
//	@Failure		400		{object}	common_model.DescriptiveError	"Missing or invalid media ID"
//	@Failure		500		{object}	common_model.DescriptiveError	"Failed to download media"
//	@Security		ApiKeyAuth
//	@Router			/media/whatsapp/download/{mediaId} [get]
func DownloadWhatsAppMedia(ctx *fiber.Ctx) error {
	mediaId := ctx.Params("mediaId")
	if mediaId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			cmn_model.NewApiError("media ID is required", errors.New("no media ID provided"), "handler").Send(),
		)
	}

	mediaInfo, err := media_service.RetrieveURL(whatsapp.WabaApi, mediaId, media_model.RetrieveInfo{})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			cmn_model.NewApiError("failed to retrieve media URL", err, "service").Send(),
		)
	}

	mediaBytes, err := media_service.Download(whatsapp.WabaApi, mediaInfo.URL)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			cmn_model.NewApiError("failed to download media", err, "service").Send(),
		)
	}

	ctx.Set("Content-Type", mediaInfo.MimeType)
	ctx.Set("Content-Disposition", "attachment; filename="+mediaId+"."+common_service.GetExtensionFromMimeType(mediaInfo.MimeType))
	ctx.Set("Content-Length", strconv.FormatInt(mediaInfo.FileSize, 10))

	return ctx.Send(mediaBytes)
}

// DownloadFromMediaInfo downloads media based on information in the request body.
//	@Summary		Download media from MediaInfo
//	@Description	Receives MediaInfo JSON, downloads the media from the provided URL, and sends it back as a file.
//	@Tags			Media
//	@Accept			json
//	@Produce		application/octet-stream
//	@Param			mediaInfo	body		media_model.MediaInfo			true	"Media Info with URL and metadata"
//	@Success		200			{file}		binary							"Downloaded media file"
//	@Failure		400			{object}	common_model.DescriptiveError	"Invalid MediaInfo"
//	@Failure		500			{object}	common_model.DescriptiveError	"Failed to download media"
//	@Security		ApiKeyAuth
//	@Router			/media/whatsapp/media-info/download [post]
func DownloadFromMediaInfo(ctx *fiber.Ctx) error {
	var mediaInfo media_model.MediaInfo
	if err := ctx.BodyParser(&mediaInfo); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			cmn_model.NewParseJsonError(err).Send(),
		)
	}

	mediaBytes, err := media_service.Download(whatsapp.WabaApi, mediaInfo.URL)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			cmn_model.NewApiError("failed to download media", err, "handler").Send(),
		)
	}

	ctx.Set("Content-Type", mediaInfo.MimeType)
	ctx.Set("Content-Disposition", "attachment; filename="+mediaInfo.Id.Id+"."+common_service.GetExtensionFromMimeType(mediaInfo.MimeType))
	ctx.Set("Content-Length", strconv.FormatInt(mediaInfo.FileSize, 10))

	return ctx.Send(mediaBytes)
}

// UploadWhatsAppMedia uploads a media file to WhatsApp.
//	@Summary		Upload media file
//	@Description	Uploads media files to WhatsApp. Files persist for up to 30 days unless deleted earlier.
//	@Tags			Media
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file							true	"Media file"
//	@Param			type	formData	string							true	"MIME type of the media file"
//	@Success		200		{object}	common_model.Id					"Media ID returned from WhatsApp"
//	@Failure		400		{object}	common_model.DescriptiveError	"Missing file or MIME type"
//	@Failure		415		{object}	common_model.DescriptiveError	"Unsupported media type"
//	@Failure		500		{object}	common_model.DescriptiveError	"Failed to upload media"
//	@Security		ApiKeyAuth
//	@Router			/media/whatsapp/upload [post]
func UploadWhatsAppMedia(ctx *fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			cmn_model.NewApiError("file is required", err, "handler").Send(),
		)
	}

	mimeType := ctx.FormValue("type")
	if mimeType == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			cmn_model.NewApiError("MIME type is required", errors.New("no type provided"), "handler").Send(),
		)
	}

	supportedMimeType, err := common_model.ParseMimeType(mimeType)
	if err != nil {
		return ctx.Status(fiber.StatusUnsupportedMediaType).JSON(
			cmn_model.NewApiError("unsupported MIME type", err, "handler").Send(),
		)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			cmn_model.NewApiError("unable to open file", err, "handler").Send(),
		)
	}
	defer file.Close()

	uploadData := media_model.Upload{
		FileName: fileHeader.Filename,
		FileData: file,
		Type:     supportedMimeType,
	}
	uploadData.SetDefault()

	mediaId, err := media_service.Upload(whatsapp.WabaApi, uploadData)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			cmn_model.NewApiError("failed to upload media", err, "handler").Send(),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(mediaId)
}
