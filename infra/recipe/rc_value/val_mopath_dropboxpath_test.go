package rc_value

import (
	"flag"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/quality/infra/qt_control"
	"testing"
)

type ValueMoPathDropboxPathRecipe struct {
	Dest mo_path.DropboxPath
}

func (z *ValueMoPathDropboxPathRecipe) Preset() {
}

func (z *ValueMoPathDropboxPathRecipe) Exec(c app_control.Control) error {
	return nil
}

func (z *ValueMoPathDropboxPathRecipe) Test(c app_control.Control) error {
	return nil
}

func TestValueMoPathDropboxPathSuccess(t *testing.T) {
	dest := "/watermint/toolbox"
	err := qt_control.WithControl(func(c app_control.Control) error {
		rcp0 := &ValueMoPathDropboxPathRecipe{}
		repo := NewRepository(rcp0)

		// Parse flags
		flg := flag.NewFlagSet("value", flag.ContinueOnError)
		repo.ApplyFlags(flg, c.UI())
		if err := flg.Parse([]string{"-dest", dest}); err != nil {
			t.Error(err)
			return err
		}

		// Apply parsed values
		rcp1 := repo.Apply()
		mod1 := rcp1.(*ValueMoPathDropboxPathRecipe)
		if !mod1.Dest.IsValid() || mod1.Dest.Path() != dest {
			t.Error(mod1)
		}

		// Spin up
		rcp2, err := repo.SpinUp(c)
		if err != nil {
			t.Error(err)
			return err
		}
		mod2 := rcp2.(*ValueMoPathDropboxPathRecipe)
		if !mod2.Dest.IsValid() || mod2.Dest.Path() != dest {
			t.Error(mod2)
		}

		if err := repo.SpinDown(c); err != nil {
			t.Error(err)
			return err
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestValueMoPathDropboxPathMissing(t *testing.T) {
	err := qt_control.WithControl(func(c app_control.Control) error {
		rcp0 := &ValueMoPathDropboxPathRecipe{}
		repo := NewRepository(rcp0)

		// Parse flags
		flg := flag.NewFlagSet("value", flag.ContinueOnError)
		repo.ApplyFlags(flg, c.UI())
		if err := flg.Parse([]string{}); err != nil {
			t.Error(err)
			return err
		}

		// Apply parsed values
		rcp1 := repo.Apply()
		mod1 := rcp1.(*ValueMoPathDropboxPathRecipe)
		if mod1.Dest.IsValid() || mod1.Dest.Path() != "" {
			t.Error(mod1)
		}

		// Spin up
		_, err := repo.SpinUp(c)
		if err == nil || err != ErrorMissingRequiredOption {
			t.Error(err)
			return err
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
