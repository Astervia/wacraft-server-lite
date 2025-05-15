package messaging_product_handler

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	_ "github.com/Astervia/wacraft-core/src/messaging-product/entity"
	messaging_product_entity "github.com/Astervia/wacraft-core/src/messaging-product/entity"
	messaging_product_model "github.com/Astervia/wacraft-core/src/messaging-product/model"
	messaging_product_service "github.com/Astervia/wacraft-server/src/messaging-product/service"
	"github.com/gofiber/fiber/v2"
)

// CreateWhatsAppContact creates a new WhatsApp contact for a messaging product.
//	@Summary		Creates a new messaging product WhatsApp contact
//	@Description	Creates and stores a new WhatsApp contact associated with a messaging product, using WhatsApp-specific product details.
//	@Tags			Messaging product contact
//	@Accept			json
//	@Produce		json
//	@Param			contact	body		messaging_product_model.CreateWhatsAppContact		true	"Contact data"
//	@Success		201		{object}	messaging_product_entity.MessagingProductContact	"Created contact"
//	@Failure		400		{object}	common_model.DescriptiveError						"Invalid request body"
//	@Failure		500		{object}	common_model.DescriptiveError						"Failed to create contact"
//	@Security		ApiKeyAuth
//	@Router			/messaging-product/contact/whatsapp [post]
func CreateWhatsAppContact(c *fiber.Ctx) error {
	// Parse the request body
	var data messaging_product_model.CreateWhatsAppContact
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewParseJsonError(err).Send(),
		)
	}

	entity, err := messaging_product_service.CreateContactForMessagingProduct(
		messaging_product_entity.MessagingProductContact{
			ContactId: data.ContactId,
			ProductDetails: &messaging_product_model.ProductDetails{
				WhatsAppProductDetails: &messaging_product_model.WhatsAppProductDetails{
					PhoneNumber: data.ProductDetails.PhoneNumber,
					WaId:        data.ProductDetails.WaId,
				},
			},
		},
		messaging_product_model.WhatsApp,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			common_model.NewApiError("unable to create whatsapp messaging product contact", err, "messaging_product_service").Send(),
		)
	}

	// Return the created user (or just a success message)
	return c.Status(fiber.StatusCreated).JSON(entity)
}
