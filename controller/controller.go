package controller

import (
	"Music/services"
	"Music/utils"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/gin-gonic/gin"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var musicService *services.MusicService

func init() {
	musicService = services.NewMusicService()
}

type MusicItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Platform string `json:"platform"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Artwork  string `json:"artwork"`
	URL      string `json:"url"`
}

// SearchMusic handles music search requests
func SearchMusic(c *gin.Context) {
	// 获取查询参数
	query := c.Query("keyword")
	searchType := c.DefaultQuery("type", "music")

	if searchType == "music" {
		results, err := musicService.SearchMusic(query)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"isEnd": true,
			"data":  results,
		})
		return
	}

	// 空类型 fallback
	c.JSON(200, gin.H{
		"isEnd": true,
		"data":  []interface{}{},
	})
}

// searchMusicFiles searches for music files matching the query
func searchMusicFiles(query string, musicRoot string) []MusicItem {
	var results []MusicItem

	err := filepath.Walk(musicRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			utils.Error("遍历路径失败:", err)
			return err
		}

		// 忽略目录
		if info.IsDir() {
			return nil
		}

		// 检查是否为音频文件
		//if !isAudioFile(info.Name()) {
		//	return nil
		//}

		// 尝试读取音频文件元数据
		metadataItem, err := extractMusicMetadata(path)
		if err != nil {
			utils.Error("无法提取元数据：", err)
			return nil
		}

		// 判断是否匹配搜索条件（大小写不敏感）
		if strings.Contains(strings.ToLower(metadataItem.Album), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(metadataItem.Artist), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(metadataItem.Name), strings.ToLower(query)) {

			// 获取文件夹名作为歌单
			dirName := filepath.Base(filepath.Dir(path))
			metadataItem.ID = filepath.Join(dirName, info.Name())
			metadataItem.URL = path
			metadataItem.Platform = "shenzaoyi"

			results = append(results, metadataItem)
			utils.Debug("匹配的音乐项: %+v", metadataItem)
		}

		return nil
	})

	if err != nil {
		utils.Error("遍历路径失败:", err)
	}

	utils.Info("总共匹配到 %d 个结果", len(results))
	return results
}

// extractMusicMetadata 提取音频文件元数据
func extractMusicMetadata(filePath string) (MusicItem, error) {
	// 打开音频文件
	file, err := os.Open(filePath)
	if err != nil {
		return MusicItem{}, err
	}
	defer file.Close()

	// 使用 tag 库读取元数据
	metadata, err := tag.ReadFrom(file)
	if err != nil {
		// 如果无法读取元数据，尝试从文件名解析
		return parseMetadataFromFileName(filePath), nil
	}

	// 创建 MusicItem 结构体
	musicItem := MusicItem{
		Name:   metadata.Title(),
		Artist: metadata.Artist(),
		Album:  metadata.Album(),
		Title:  metadata.Title(),
	}

	// 尝试获取专辑封面
	picture := metadata.Picture()
	if picture != nil {
		// 保存专辑封面到临时文件
		artworkPath := filepath.Join(os.TempDir(), fmt.Sprintf("%s_artwork.jpg", musicItem.Name))
		err = os.WriteFile(artworkPath, picture.Data, 0644)
		if err == nil {
			musicItem.Artwork = artworkPath
		}
	}

	return musicItem, nil
}

// parseMetadataFromFileName 从文件名解析元数据
func parseMetadataFromFileName(filePath string) MusicItem {
	fileName := filepath.Base(filePath)

	// 移除文件扩展名
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// 尝试按常见分隔符拆分
	separators := []string{" - ", "-", "_"}
	var parts []string

	for _, sep := range separators {
		parts = strings.Split(fileName, sep)
		if len(parts) >= 2 {
			break
		}
	}

	musicItem := MusicItem{
		Name:  fileName,
		Title: fileName,
	}

	// 根据拆分结果填充信息
	switch len(parts) {
	case 3: // 专辑 - 歌手 - 歌曲
		musicItem.Album = strings.TrimSpace(parts[0])
		musicItem.Artist = strings.TrimSpace(parts[1])
		musicItem.Name = strings.TrimSpace(parts[2])
		musicItem.Title = strings.TrimSpace(parts[2])
	case 2: // 歌手 - 歌曲
		musicItem.Artist = strings.TrimSpace(parts[0])
		musicItem.Name = strings.TrimSpace(parts[1])
		musicItem.Title = strings.TrimSpace(parts[1])
	}

	return musicItem
}

// isAudioFile 检查是否为音频文件
func isAudioFile(filename string) bool {
	audioExtensions := []string{".mp3", ".flac", ".wav", ".m4a", ".ogg", ".aac"}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, audioExt := range audioExtensions {
		if ext == audioExt {
			return true
		}
	}

	return false
}

// PlayMusic handles music play requests
func PlayMusic(c *gin.Context) {
	// URL 解码
	id, err := url.QueryUnescape(c.Query("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	musicRoot := "/home/ftpuser/Music"
	musicPath := filepath.Join(musicRoot, id)

	// 如果直接路径不存在，尝试模糊匹配
	if _, err := os.Stat(musicPath); os.IsNotExist(err) {
		musicPath = findMusicFileWithVariants(musicRoot, id)
		if musicPath == "" {
			utils.Warn("音乐文件未找到: %s", id)
			c.JSON(404, gin.H{"error": "music not found"})
			return
		}
	}

	// Set proper file headers
	c.Header("Content-Type", "audio/mpeg")
	c.Header("Accept-Ranges", "bytes")

	// Send the music file to the client
	utils.Info("正在提供音乐文件: %s", musicPath)
	c.File(musicPath)
}

// findMusicFileWithVariants 尝试查找文件的多种变体
func findMusicFileWithVariants(musicRoot, id string) string {
	// 提取可能的歌曲名
	parts := strings.Split(id, "-")
	if len(parts) < 3 {
		return ""
	}
	songName := parts[len(parts)-2]

	// 可能的搜索变体
	variants := []string{
		songName,
		strings.ReplaceAll(songName, " ", ""),
		strings.ReplaceAll(songName, " ", "-"),
		id,
	}

	for _, variant := range variants {
		pattern := filepath.Join(musicRoot, "*", fmt.Sprintf("*%s*", variant))
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			return matches[0]
		}
	}

	return ""
}

type AlbumRequest struct {
	URL string `json:"url"`
}

func GetAlbumMusics(c *gin.Context) {
	// 获取协议（http/https）
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// 获取完整的基础 URL
	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

	var req AlbumRequest

	// 使用 ShouldBindJSON 绑定 JSON 请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error("参数绑定错误:", err)
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	keyword := req.URL
	utils.Info("接收到的搜索关键词: %s", keyword)

	if keyword == "" {
		utils.Warn("关键词为空")
		c.JSON(400, gin.H{"error": "关键词不能为空"})
		return
	}

	musicRoot := "/home/ftpuser/Music"
	utils.Debug("音乐根目录: %s", musicRoot)

	// 读取一级目录
	dirs, err := os.ReadDir(musicRoot)
	if err != nil {
		utils.Error("读取目录失败:", err)
		c.JSON(500, gin.H{"error": "读取目录失败"})
		return
	}

	var matchedDir string

	// 查找包含关键词的目录
	utils.Debug("开始遍历目录")
	for _, dir := range dirs {
		if dir.IsDir() {
			utils.Debug("检查目录: %s", dir.Name())
			if strings.Contains(strings.ToLower(dir.Name()), strings.ToLower(keyword)) {
				matchedDir = dir.Name()
				utils.Info("找到匹配目录: %s", matchedDir)
				break
			}
		}
	}

	if matchedDir == "" {
		utils.Warn("未找到匹配的专辑")
		c.JSON(404, gin.H{"error": "未找到匹配的专辑"})
		return
	}

	// 读取匹配目录下的音乐文件
	albumPath := filepath.Join(musicRoot, matchedDir)
	utils.Debug("专辑路径: %s", albumPath)

	files, err := os.ReadDir(albumPath)
	if err != nil {
		utils.Error("读取专辑失败:", err)
		c.JSON(500, gin.H{"error": "读取专辑失败"})
		return
	}

	var musicList []map[string]string
	utils.Debug("开始遍历音乐文件")
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(albumPath, file.Name())

			// 尝试提取元数据
			metadata, err := extractMusicMetadata(filePath)
			if err != nil {
				utils.Error("无法提取元数据：", err)
				continue
			}

			// 构建完整的播放 URL
			playURL := fmt.Sprintf("%s/music/v1/play?id=%s", baseURL, url.QueryEscape(filepath.Join(matchedDir, file.Name())))

			musicItem := map[string]string{
				"id":      file.Name(),
				"title":   metadata.Title,
				"artist":  metadata.Artist,
				"album":   metadata.Album,
				"artwork": metadata.Artwork,
				"url":     playURL,
			}

			musicList = append(musicList, musicItem)
		}
	}

	utils.Info("总共找到 %d 首音乐", len(musicList))
	c.JSON(200, gin.H{
		"data": musicList,
	})
}

func GetAlbumList(c *gin.Context) {
	musicRoot := "/home/ftpuser/Music"

	// 读取一级目录
	dirs, err := os.ReadDir(musicRoot)
	if err != nil {
		utils.Error("读取目录失败:", err)
		c.JSON(500, gin.H{"error": "读取目录失败"})
		return
	}

	var albumList []string
	for _, dir := range dirs {
		if dir.IsDir() {
			albumList = append(albumList, dir.Name())
			utils.Debug("找到专辑目录: %s", dir.Name())
		}
	}
	utils.Info("总共找到 %d 个专辑", len(albumList))

	c.JSON(200, gin.H{
		"data": albumList,
	})
}

type MusicInfo struct {
	Singer   string
	Album    string
	Name     string
	Cover    string
	Location string
}

// upload music, temp
//func UploadMusic(c *gin.Context) {
//	cosClient, err := tengcent_cos.InitClient()
//	if err != nil {
//		fmt.Println("error")
//	}
//	// upload music file
//	cosClient.Upload("")
//	// fill meat data
//	musicInfo := MusicInfo{
//		Singer:   "",
//		Album:    "",
//		Name:     "",
//		Cover:    "",
//		Location: "",
//	}
//
//	// save in database
//}
