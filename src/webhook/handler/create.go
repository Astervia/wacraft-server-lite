package webhook_handler

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	"github.com/Astervia/wacraft-core/src/repository"
	webhook_entity "github.com/Astervia/wacraft-core/src/webhook/entity"
	webhook_model "github.com/Astervia/wacraft-core/src/webhook/model"
	"github.com/Astervia/wacraft-server/src/database"
	"github.com/gofiber/fiber/v2"
)

// @Summary		Create a new webhook
// @Description	Creates a new webhook with the specified URL, authorization, and event type
// @Tags			Webhook
// @Accept			json
// @Produce		json
// @Param			webhook	body		webhook_model.CreateWebhook		true	"Webhook data"
// @Success		201		{object}	webhook_entity.Webhook			"Created webhook"
// @Failure		400		{object}	common_model.DescriptiveError	"Invalid request body"
// @Failure		500		{object}	common_model.DescriptiveError	"Internal server error"
// @Router			/webhook [post]
// @Security		ApiKeyAuth
func CreateWebhook(c *fiber.Ctx) error {
	var newWebhook webhook_model.CreateWebhook
	if err := c.BodyParser(&newWebhook); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewParseJsonError(err).Send(),
		)
	}

	webhook, err := repository.Create(
		webhook_entity.Webhook{
			Url:           newWebhook.Url,
			Authorization: newWebhook.Authorization,
			HttpMethod:    newWebhook.HttpMethod,
			Timeout:       newWebhook.Timeout,
			Event:         newWebhook.Event,
		}, database.DB,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			common_model.NewApiError("unable to create webhook", err, "repository").Send(),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(webhook)
}
