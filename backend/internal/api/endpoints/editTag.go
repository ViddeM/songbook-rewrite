package endpoints

import (
	"errors"
	"github.com/ericlp/songbook/backend/internal/common"
	"github.com/ericlp/songbook/backend/internal/db/queries"
	"github.com/ericlp/songbook/backend/internal/db/tables"
	"github.com/ericlp/songbook/backend/internal/process"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func EditTag(c *gin.Context) {
	tag, err := validateTag(c)
	if err != nil {
		log.Printf("Failed to validate edit tag json: %v\n", err)
		c.JSON(http.StatusBadRequest, common.Error(common.ResponseInvalidJson))
		return
	}

	oldTag, err := validateTagId(c)
	if err != nil {
		log.Printf("Failed to validate tag id: %v\n", err)
		return
	}

	err = validateOwnerAuthorized(c, oldTag.OwnedBy)
	if err != nil {
		log.Printf("User not authorized to edit recipe: %v\n")
		c.JSON(http.StatusForbidden, common.Error(common.ResponseIncorrectUser))
		return
	}

	err = process.EditTag(oldTag, tag)
	if err != nil {
		log.Printf("Failed to edit tag: %v\n", err)
		if errors.Is(err, common.ErrNameTaken) {
			c.JSON(
				http.StatusUnprocessableEntity,
				common.Error(common.ResponseTagNameTaken),
			)
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			common.Error(common.ResponseFailedToEditTag),
		)
		return
	}

	c.JSON(http.StatusOK, common.Success("tag edited"))
}

func validateTagId(c *gin.Context) (*tables.Tag, error) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			common.Error(common.ResponseMalformedTagId),
		)
		return nil, err
	}

	tag, err := queries.GetTagById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.Error(common.ResponseTagNotFound))
		return nil, err
	}

	return tag, nil
}
