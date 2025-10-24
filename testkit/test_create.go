package testkit

import (
	"github.com/trateou/gorm-cache/cache"
	. "github.com/smartystreets/goconvey/convey"

	"gorm.io/gorm"
)

func testSearchCreate(cache *cache.Gorm2Cache, db *gorm.DB) {
	err := cache.ResetCache()
	So(err, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	var model = new(TestModel)

	result := db.Where("id = ?", 1).First(model)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)
	So(model.ID, ShouldEqual, 1)

	result = db.Create(&TestModel{})
	So(result.Error, ShouldBeNil)

	result = db.Where("id = ?", 1).First(model)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)
}
