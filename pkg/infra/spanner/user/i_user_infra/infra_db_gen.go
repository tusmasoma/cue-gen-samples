// Code generated by cue_gen. DO NOT EDIT.
package i_user_infra

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/pkg/errors"

	"github.com/tusmasoma/cue-gen-samples/pkg/infra/spanner/model/user/i_user"
)

// IFInfra - インフラ層のインターフェース
type IFInfra interface {
	Get(ctx context.Context,
		userId string,
	) (*i_user.Model, error)

	TryGet(ctx context.Context,
		userId string,
	) (*i_user.Model, bool, error)

	// FindByConditions(ctx context.Context, conds ...spanner.Condition) (i_user.Models, error)

	Insert(ctx context.Context, model *i_user.Model) error
	Update(ctx context.Context, model *i_user.Model) error
	Delete(ctx context.Context, model *i_user.Model) error
	Save(ctx context.Context, model *i_user.Model) error
}

// NewInfra - 新規作成
func NewInfra(client *spanner.Client) *Infra {
	return &Infra{client: client}
}

// Infra - インフラ層の実装
type Infra struct {
	client *spanner.Client
}

var _ IFInfra = (*Infra)(nil)

// Get - 取得
func (i *Infra) Get(ctx context.Context,
	userId string,
) (*i_user.Model, error) {
	row, err := i.client.Single().ReadRow(ctx, i_user.TableName, spanner.Key{
		userId,
	}, i_user.Columns)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get data from spanner")
	}

	var model i_user.Model
	if err := model.Bind(row); err != nil {
		return nil, errors.Wrap(err, "failed to bind row data")
	}

	return &model, nil
}

// TryGet - 取得（エラー時に false を返す）
func (i *Infra) TryGet(ctx context.Context,
	userId string,
) (*i_user.Model, bool, error) {
	model, err := i.Get(ctx, userId)
	if err != nil {
		if spanner.ErrCode(err) == spanner.ErrRowNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return model, true, nil
}

// Insert - データを挿入
func (i *Infra) Insert(ctx context.Context, model *i_user.Model) error {
	mutation := model.InsertMutation()
	_, err := i.client.Apply(ctx, []*spanner.Mutation{mutation})
	return errors.Wrap(err, "failed to insert data into spanner")
}

// Update - データを更新
func (i *Infra) Update(ctx context.Context, model *i_user.Model) error {
	mutation := model.UpdateMutation()
	_, err := i.client.Apply(ctx, []*spanner.Mutation{mutation})
	return errors.Wrap(err, "failed to update data in spanner")
}

// Delete - データを削除
func (i *Infra) Delete(ctx context.Context, model *i_user.Model) error {
	mutation := model.DeleteMutation()
	_, err := i.client.Apply(ctx, []*spanner.Mutation{mutation})
	return errors.Wrap(err, "failed to delete data from spanner")
}

// Save - データを保存（存在しなければ Insert、存在すれば Update）
func (i *Infra) Save(ctx context.Context, model *i_user.Model) error {
	_, exists, err := i.TryGet(ctx, model.UserId)
	if err != nil {
		return err
	}
	if exists {
		return i.Update(ctx, model)
	}
	return i.Insert(ctx, model)
}

// FindByConditions - 条件検索
/* func (i *Infra) FindByConditions(ctx context.Context, conds ...spanner.Condition) (i_user.Models, error) {
	stmt := spanner.Statement{
		SQL: "SELECT * FROM " + i_user.TableName + " WHERE " + spanner.BuildWhereClause(conds),
		Params: spanner.BuildParams(conds),
	}

	iter := i.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var models i_user.Models
	for {
		row, err := iter.Next()
		if err == spanner.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to query spanner")
		}

		var model i_user.Model
		if err := model.Bind(row); err != nil {
			return nil, errors.Wrap(err, "failed to bind row data")
		}
		models = append(models, &model)
	}

	return models, nil
} */
