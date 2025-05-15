package message_handler

import (
	"net/url"

	common_model "github.com/Astervia/wacraft-core/src/common/model"
	message_entity "github.com/Astervia/wacraft-core/src/message/entity"
	message_model "github.com/Astervia/wacraft-core/src/message/model"
	message_service "github.com/Astervia/wacraft-server/src/message/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetConversation returns paginated messages from a specific messaging product contact.
//	@Summary		Get conversation messages
//	@Description	Returns a paginated list of messages sent or received by the specified messaging product contact.
//	@Tags			Message conversation
//	@Accept			json
//	@Produce		json
//	@Param			message						query		message_model.QueryPaginated	true	"Pagination and filter parameters"
//	@Param			messagingProductContactId	path		string							true	"Messaging product contact ID"
//	@Success		200							{array}		message_entity.Message			"Conversation messages"
//	@Failure		400							{object}	common_model.DescriptiveError	"Invalid query or ID"
//	@Failure		500							{object}	common_model.DescriptiveError	"Failed to retrieve messages"
//	@Security		ApiKeyAuth
//	@Router			/message/conversation/messaging-product-contact/{messagingProductContactId} [get]
func GetConversation(c *fiber.Ctx) error {
	mpcId, err := uuid.Parse(c.Params("messagingProductContactId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewApiError("unable to parse messaging product contact id string to UUID", err, "github.com/google/uuid"),
		)
	}

	query := new(message_model.QueryPaginated)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewParseJsonError(err).Send(),
		)
	}

	messages, err := message_service.GetConversation(
		mpcId,
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
		"",
		nil,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			common_model.NewApiError("unable to get conversation messages", err, "message_service"),
		)
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}

// GetConversations returns the latest message in each conversation.
//	@Summary		Get conversations
//	@Description	Returns a paginated list of the latest messages per conversation, enriched with contact information.
//	@Tags			Message conversation
//	@Accept			json
//	@Produce		json
//	@Param			message	query		message_model.QueryPaginated	true	"Pagination and filter parameters"
//	@Success		200		{array}		message_entity.Message			"Latest messages per conversation"
//	@Failure		400		{object}	common_model.DescriptiveError	"Invalid query"
//	@Failure		500		{object}	common_model.DescriptiveError	"Failed to retrieve conversations"
//	@Security		ApiKeyAuth
//	@Router			/message/conversation [get]
func GetConversations(c *fiber.Ctx) error {
	query := new(message_model.QueryPaginated)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewParseJsonError(err),
		)
	}

	messages, err := message_service.GetLatestMessagesForEachUser(
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
			common_model.NewApiError("unable to get conversations", err, "message_service"),
		)
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}

// CountConversationsByMessagingProductContact counts the number of messages exchanged with a specific contact.
//	@Summary		Counts conversations by messaging product contact
//	@Description	Counts messages exchanged with the specified messaging product contact based on filters.
//	@Tags			Message conversation
//	@Accept			json
//	@Produce		json
//	@Param			message						query		message_model.QueryPaginated	true	"Filter parameters"
//	@Param			messagingProductContactId	path		string							true	"Messaging product contact ID"
//	@Success		200							{integer}	int								"Count of messages"
//	@Failure		400							{object}	common_model.DescriptiveError	"Invalid query or ID"
//	@Failure		500							{object}	common_model.DescriptiveError	"Failed to count messages"
//	@Security		ApiKeyAuth
//	@Router			/message/conversation/count/messaging-product-contact/{messagingProductContactId} [get]
func CountConversationsByMessagingProductContact(c *fiber.Ctx) error {
	mpcId, err := uuid.Parse(c.Params("messagingProductContactId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewApiError("unable to parse messaging product contact id string to UUID", err, "github.com/google/uuid"),
		)
	}

	query := new(message_model.Query)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewParseJsonError(err),
		)
	}

	messages, err := message_service.CountConversations(
		mpcId,
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
		&query.DateOrder,
		&query.DateWhereWithDeletedAt,
		"",
		nil,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			common_model.NewApiError("unable to count conversations", err, "message_service"),
		)
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}

// CountDistinctConversations returns the number of unique conversations.
//	@Summary		Counts conversations
//	@Description	Counts distinct conversations based on the provided filters.
//	@Tags			Message conversation
//	@Accept			json
//	@Produce		json
//	@Param			message	query		message_model.QueryPaginated	true	"Filter parameters"
//	@Success		200		{integer}	int								"Count of distinct conversations"
//	@Failure		400		{object}	common_model.DescriptiveError	"Invalid query"
//	@Failure		500		{object}	common_model.DescriptiveError	"Failed to count conversations"
//	@Security		ApiKeyAuth
//	@Router			/message/conversation/count [get]
func CountDistinctConversations(c *fiber.Ctx) error {
	query := new(message_model.Query)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewParseJsonError(err),
		)
	}

	messages, err := message_service.CountDistinctConversations(
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
		&query.DateOrder,
		&query.DateWhereWithDeletedAt,
		"",
		nil,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			common_model.NewApiError("unable to count distinct conversations", err, "message_service"),
		)
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}

// ConversationContentLikeByMessagingProductContact returns messages matching text content.
//	@Summary		Count conversation messages by content and messaging product contact ID
//	@Description	Returns messages filtered by a "like" match on sender/receiver/product data fields.
//	@Tags			Message conversation
//	@Accept			json
//	@Produce		json
//	@Param			message						query		message_model.QueryPaginated	true	"Filter parameters"
//	@Param			messagingProductContactId	path		string							true	"Messaging product contact ID"
//	@Param			likeText					path		string							true	"Substring to match against sender/receiver/product data"
//	@Success		200							{array}		message_entity.Message			"Filtered conversation messages"
//	@Failure		400							{object}	common_model.DescriptiveError	"Invalid ID, query, or likeText"
//	@Failure		500							{object}	common_model.DescriptiveError	"Failed to retrieve messages"
//	@Security		ApiKeyAuth
//	@Router			/message/conversation/messaging-product-contact/{messagingProductContactId}/content/like/{likeText} [get]
func ConversationContentLikeByMessagingProductContact(c *fiber.Ctx) error {
	// Parse the messagingProductContactId from the path
	mpcId, err := uuid.Parse(c.Params("messagingProductContactId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewApiError("unable to parse messaging product contact id string to UUID", err, "github.com/google/uuid").Send(),
		)
	}

	// Parse and decode the likeText from the path
	encodedText := c.Params("likeText")
	decodedText, err := url.QueryUnescape(encodedText)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewApiError("unable to decode likeText", err, "net/url").Send(),
		)
	}

	// Parse query parameters into the QueryPaginated model
	query := new(message_model.QueryPaginated)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewParseJsonError(err).Send(),
		)
	}

	// Call the ConversationContentLike service
	messages, err := message_service.ConversationContentLike(
		mpcId,
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
		nil, // Pass a nil *gorm.DB to use the default connection
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			common_model.NewApiError("unable to get conversation messages by content", err, "message_service").Send(),
		)
	}

	// Return the filtered messages as JSON
	return c.Status(fiber.StatusOK).JSON(messages)
}

// CountConversationContentLike counts messages matching a likeText.
//	@Summary		Counts conversation messages by content
//	@Description	Counts messages from a messaging product contact matching a "like" pattern on sender/receiver/product data.
//	@Tags			Message conversation
//	@Accept			json
//	@Produce		json
//	@Param			message						query		message_model.QueryPaginated	true	"Filter parameters"
//	@Param			messagingProductContactId	path		string							true	"Messaging product contact ID"
//	@Param			likeText					path		string							true	"Substring to match against sender/receiver/product data"
//	@Success		200							{integer}	int								"Count of matched messages"
//	@Failure		400							{object}	common_model.DescriptiveError	"Invalid query or likeText"
//	@Failure		500							{object}	common_model.DescriptiveError	"Failed to count messages"
//	@Security		ApiKeyAuth
//	@Router			/message/conversation/count/messaging-product-contact/{messagingProductContactId}/content/like/{likeText} [get]
func CountConversationContentLike(c *fiber.Ctx) error {
	// Parse the messagingProductContactId from the path
	mpcId, err := uuid.Parse(c.Params("messagingProductContactId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewApiError("unable to parse messaging product contact id string to UUID", err, "github.com/google/uuid").Send(),
		)
	}

	// Parse and decode the likeText from the path
	encodedText := c.Params("likeText")
	decodedText, err := url.QueryUnescape(encodedText)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewApiError("unable to decode likeText", err, "net/url").Send(),
		)
	}

	// Parse query parameters into the QueryPaginated model
	query := new(message_model.Query)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			common_model.NewParseJsonError(err).Send(),
		)
	}

	// Call the ConversationContentLike service
	count, err := message_service.CountConversationContentLike(
		mpcId,
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
		&query.DateOrder,
		&query.DateWhereWithDeletedAt,
		nil, // Pass a nil *gorm.DB to use the default connection
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			common_model.NewApiError("unable to count conversation messages by content", err, "message_service").Send(),
		)
	}

	// Return the filtered messages as JSON
	return c.Status(fiber.StatusOK).JSON(count)
}
