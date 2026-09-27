package controllers

import (
	"exchangeapp/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 同一个用户只能点赞一次：用 Redis Set 记录赞过的用户，天然去重。
func LikeArticle(ctx *gin.Context) {
	articleID := ctx.Param("id")
	username := ctx.GetString("username")
	if username == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	likeKey := "article:" + articleID + ":liked_users"

	added, err := global.RedisDB.SAdd(likeKey, username).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if added == 0 {
		ctx.JSON(http.StatusConflict, gin.H{"message": "Already liked this article"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully liked the article"})
}

func GetArticleLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")
	likeKey := "article:" + articleID + ":liked_users"

	count, err := global.RedisDB.SCard(likeKey).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"likes": count})
}
