package handlers

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/yosephbuitrago/garagesale/business/auth"
	"github.com/yosephbuitrago/garagesale/business/data/product"
	"github.com/yosephbuitrago/garagesale/foundation/web"
)

type productGroup struct {
	product product.Product
}

func (pg productGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	products, err := pg.product.Query(ctx, v.TraceID)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, products, http.StatusOK)
}

func (pg productGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("claims missing from context")
	}

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return errors.Wrap(err, "decoding new product")
	}

	prod, err := pg.product.Create(ctx, v.TraceID, claims, np, v.Now)
	if err != nil {
		return errors.Wrapf(err, "creating new product %+v", np)
	}

	return web.Respond(ctx, w, prod, http.StatusCreated)
}
