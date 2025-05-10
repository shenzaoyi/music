package services

import (
	"Music/config"
	"Music/models"
	"Music/repositories"
	"Music/tengcent_cos"
	"errors"
	"fmt"
	"strconv"
)

type MusicService struct {
	repo *repositories.MusicRepository
}

// 创建一个新的 MusicService 实例
func NewMusicService() *MusicService {
	return &MusicService{
		repo: &repositories.MusicRepository{},
	}
}

// 创建音乐记录（带去重判断）
func (s *MusicService) CreateMusic(info *models.MusicInfo, filepath string) error {
	// 检查是否已存在同名歌曲（可选逻辑）
	exists, err := s.repo.ExistsByFields(info.Name, info.Album, info.Singer)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("音乐记录已存在")
	}
	// 存储到数据库
	id, err := s.repo.Create(info)
	if err != nil {
		return err
	}
	// 存储文件内容到云存储，得到url
	client, err := tengcent_cos.InitClient()
	if err != nil {
		return errors.New("New tengcent client error")
	}
	location := client.Upload(strconv.Itoa(int(info.ID)), filepath)
	fmt.Println(location)
	info.Location = location
	// 跟新数据库
	updates := map[string]interface{}{
		"Location": location,
	}
	err = s.repo.Update(id, updates)
	if err != nil {
		return err
	}
	return nil
}

// 根据 ID 获取音乐记录
func (s *MusicService) GetMusic(id uint) (*models.MusicInfo, error) {
	return s.repo.GetByID(id)
}

// 更新音乐信息
func (s *MusicService) UpdateMusic(id uint, updates map[string]interface{}) error {
	return s.repo.Update(id, updates)
}

// 删除音乐记录
func (s *MusicService) DeleteMusic(id uint) error {
	return s.repo.Delete(id)
}

// 搜索结果结构体（前端需要的结构）
//
//	type SearchResult struct {
//		ID   string `json:"id"`
//		Name string `json:"name"`
//		URL  string `json:"url,omitempty"` // 最终由 handler 填充
//	}
type SearchResult struct {
	ID       string `json:"id"`
	Name     string `json:"name"`  // 临时存储Location, 也就是腾讯云对象存储的路径
	Title    string `json:"title"` // title 是歌曲名，前端忘记怎么写的了
	Platform string `json:"platform"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Artwork  string `json:"artwork"`
	URL      string `json:"url"`
}

func (s *MusicService) SearchMusic(keyword string) ([]SearchResult, error) {
	// 模糊查询
	musics, err := s.repo.SearchByKeyword(keyword)
	if err != nil {
		return nil, err
	}
	// 转换为 SearchResult 结构
	var results []SearchResult
	for _, m := range musics {
		results = append(results, SearchResult{
			ID:       strconv.Itoa(int(m.ID)),
			Name:     m.Location,
			Title:    m.Name,
			Platform: "shenzaoyi",
			Artist:   m.Singer,
			Album:    m.Album,
			Artwork:  m.Cover,
			URL:      config.PLAYBASEURL + strconv.Itoa(int(m.ID)),
		})
	}
	return results, nil
}
func (s *MusicService) GetByID(id uint) (*models.MusicInfo, error) {
	return s.repo.GetByID(id)
}
