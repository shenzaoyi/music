package repositories

import "Music/models"

type MusicRepository struct{}

func (r *MusicRepository) Create(info *models.MusicInfo) (uint, error) {
	result := models.DB.Create(info)
	return info.ID, result.Error
}

func (r *MusicRepository) GetByID(id uint) (*models.MusicInfo, error) {
	var music models.MusicInfo
	err := models.DB.First(&music, id).Error
	return &music, err
}

func (r *MusicRepository) Update(id uint, updates map[string]interface{}) error {
	return models.DB.Model(&models.MusicInfo{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MusicRepository) Delete(id uint) error {
	return models.DB.Delete(&models.MusicInfo{}, id).Error
}

// 模糊搜索函数
func (r *MusicRepository) SearchByKeyword(keyword string) ([]models.MusicInfo, error) {
	var results []models.MusicInfo
	err := models.DB.Where("name LIKE ? OR singer LIKE ? OR album LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Find(&results).Error
	return results, err
}
