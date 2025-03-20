package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/oden/internal/config"
	"github.com/yourusername/oden/internal/db"
	"github.com/yourusername/oden/internal/storage"
)

// RegisterHandlers registers all API handlers
func RegisterHandlers(router *gin.Engine, db *db.DB, storage *storage.Client, cfg *config.Config) {
	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Auth routes
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", registerHandler)
			authRoutes.POST("/login", loginHandler)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(authMiddleware(cfg))
		{
			// Heroes routes
			heroesRoutes := protected.Group("/heroes")
			{
				heroesRoutes.GET("/list", listHeroesHandler)
				heroesRoutes.POST("/summon", summonHeroHandler)
			}

			// Team routes
			teamRoutes := protected.Group("/team")
			{
				teamRoutes.GET("/get", getTeamHandler)
				teamRoutes.POST("/save", saveTeamHandler)
			}

			// Battle routes
			battleRoutes := protected.Group("/battle")
			{
				battleRoutes.POST("/start", startBattleHandler)
			}

			// Idle routes
			idleRoutes := protected.Group("/idle")
			{
				idleRoutes.GET("/rewards", getIdleRewardsHandler)
				idleRoutes.POST("/claim", claimIdleRewardsHandler)
			}

			// Items routes
			itemsRoutes := protected.Group("/items")
			{
				itemsRoutes.GET("/list", listItemsHandler)
				itemsRoutes.POST("/use", useItemHandler)
				itemsRoutes.POST("/equip", equipItemHandler)
				itemsRoutes.POST("/unequip", unequipItemHandler)
			}

			// Missions routes
			missionsRoutes := protected.Group("/missions")
			{
				missionsRoutes.GET("/list", listMissionsHandler)
				missionsRoutes.POST("/claim", claimMissionRewardHandler)
			}

			// Gacha routes
			gachaRoutes := protected.Group("/gacha")
			{
				gachaRoutes.GET("/banners", listBannersHandler)
				gachaRoutes.POST("/summon", summonGachaHandler)
				gachaRoutes.GET("/rates", getBannerRatesHandler)
			}
		}
	}
} 