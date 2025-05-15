package message_handler

import (
	"net/url"

	common_model "github.com/Astervia/wacraft-core/src/common/model"
	message_entity "github.com/Astervia/wacraft-core/src/message/entity"
	message_model "github.com/Astervia/wacraft-core/src/message/model"
	message_service "github.com/Astervia/wacraft-server/src/message/service"
	"github.com/gofiber/fiber/v2"
)

// ContentLike returns messages where text matches sender, receiver, or product data fields.
//	@Summary		Queries message content like text paginated
//	@Description	Uses regex (~) to match the given text against `sender_data`, `receiver_data`, and `product_data` fields.
//	@Tags			Message
//	@Accept			json
//	@Produce		json
//	@Param			message		query		message_model.QueryPaginated	true	"Pagination and filter parameters"
//	@Param			likeText	path		string							true	"Text to apply like operator"
//	@Success		200			{array}		message_entity.Message			"List of matched messages"
//	@Failure		400			{object}	common_model.DescriptiveError	"Invalid likeText or query"
//	@Failure		500			{object}	common_model.DescriptiveError	"Failed to query messages"
//	@Security		ApiKeyAuth
//	@Router			/message/content/like/{likeText} [get]
func ContentLike(c *fiber.Ctx) error {
	encodedText := c.Params("likeText")
	decodedText, err := url.QueryUnescape(encodedText)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewApiError("unable to decode likeText", err, "net/url").Send(),
		)
	}

	query := new(message_model.QueryPaginated)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewParseJsonError(err).Send(),
		)
	}

	messages, err := message_service.ContentLike(
		decodedText,
		message_entity.Message{
			MessageFields: message_model.MessageFields{
				FromId:             query.FromId,
				ToId:               query.ToId,
				MessagingProductId: query.MessagingProductId,
				AuditWithDeleted: common_model.AuditWithDeleted{
					Audit: common_model.Audit{
						Id: query.Id,
					},
				},
			},
		},
		&query.Paginate,
		&query.DateOrder,
		&query.DateWhereWithDeletedAt,
		nil,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			common_model.NewApiError("unable to get messages", err, "message_service").Send(),
		)
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}

// ContentKeyLike returns messages matching text in a specific field.
//	@Summary		Queries message content like text paginated
//	@Description	Uses regex (~) to match the given text on a dynamic key field. The fields `from` and `to` are populated in the result.
//	@Tags			Message
//	@Accept			json
//	@Produce		json
//	@Param			message		query		message_model.QueryPaginated	true	"Pagination and filter parameters"
//	@Param			keyName		path		string							true	"Field name to apply the like operator"
//	@Param			likeText	path		string							true	"Text to apply like operator"
//	@Success		200			{array}		message_entity.Message			"List of matched messages"
//	@Failure		400			{object}	common_model.DescriptiveError	"Invalid keyName, likeText, or query"
//	@Failure		500			{object}	common_model.DescriptiveError	"Failed to query messages"
//	@Security		ApiKeyAuth
//	@Router			/message/content/{keyName}/like/{likeText} [get]
func ContentKeyLike(c *fiber.Ctx) error {
	encodedText := c.Params("likeText")
	decodedText, err := url.QueryUnescape(encodedText)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewApiError("unable to decode likeText", err, "net/url").Send(),
		)
	}

	encodedKey := c.Params("keyName")
	decodedKey, err := url.QueryUnescape(encodedKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewApiError("unable to decode keyName", err, "net/url").Send(),
		)
	}

	query := new(message_model.QueryPaginated)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewParseJsonError(err).Send(),
		)
	}

	messages, err := message_service.ContentKeyLike(
		decodedText,
		decodedKey,
		message_entity.Message{
			MessageFields: message_model.MessageFields{
				FromId:             query.FromId,
				ToId:               query.ToId,
				MessagingProductId: query.MessagingProductId,
				AuditWithDeleted: common_model.AuditWithDeleted{
					Audit: common_model.Audit{
						Id: query.Id,
					},
				},
			},
		},
		&query.Paginate,
		&query.DateOrder,
		&query.DateWhereWithDeletedAt,
		nil,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			common_model.NewApiError("unable to get messages", err, "message_service").Send(),
		)
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}
