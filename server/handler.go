package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"restful-ethereum-validator/service"
	"strconv"
)

type APIHandler struct {
	ethereumService service.EthereumService
	logger          *logrus.Logger
}

func NewAPIHandler(logger *logrus.Logger, ethereumService service.EthereumService) *APIHandler {
	return &APIHandler{
		ethereumService: ethereumService,
		logger:          logger,
	}
}

// GetBlockReward returns total block reward in gwei for a given 'slot'
func (h *APIHandler) GetBlockReward(c *gin.Context) {
	slotStr := c.Param("slot")
	slot, err := strconv.Atoi(slotStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot number"})
		return
	}
	currentSlot, err := h.ethereumService.GetCurrentSlot(c)
	if err != nil {
		h.logger.Errorf("Failed to get current slot: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current slot"})
		return
	}

	if slot > int(currentSlot) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requested slot is in the future"})
		return
	}

	rewardStatus, rewardAmount, err := h.ethereumService.GetBlockReward(c, uint64(slot))
	if err != nil {
		h.logger.Errorf("Failed to get block reward: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get block reward"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": rewardStatus, "reward": rewardAmount})
}

// GetSyncDuties returns a list of validator's public keys that have sync committee duties for a given 'slot'
func (h *APIHandler) GetSyncDuties(c *gin.Context) {
	slotStr := c.Param("slot")
	slot, err := strconv.Atoi(slotStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot number"})
		return
	}

	currentSlot, err := h.ethereumService.GetCurrentSlot(c)
	if err != nil {
		h.logger.Errorf("Failed to get current slot: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current slot"})
		return
	}

	if slot > int(currentSlot) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requested slot is in the future"})
		return
	}

	syncDuties, err := h.ethereumService.GetSyncDuties(c, uint64(slot))
	if err != nil {
		h.logger.Errorf("Failed to get sync duties: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sync duties"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"validators": syncDuties})
}
