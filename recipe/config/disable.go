package config

import (
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/control/app_control_launcher"
	"github.com/watermint/toolbox/infra/control/app_feature"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/ui/app_msg"
)

type Disable struct {
	Key                         string
	ErrorInvalidKey             app_msg.Message
	ErrorUnableToDisableFeature app_msg.Message
	InfoOptOut                  app_msg.Message
}

func (z *Disable) Preset() {
}

func (z *Disable) Exec(c app_control.Control) error {
	l := c.Log()
	ui := c.UI()
	cl, ok := c.(app_control_launcher.ControlLauncher)
	if !ok {
		l.Debug("Catalogue is not available")
		return ErrorCatalogueIsNotAvailable
	}
	features := cl.Catalogue().Features()
	if c.Feature().IsTest() {
		features = append(features, &SampleFeature{})
	}
	var feature app_feature.OptIn = nil
	for _, f := range features {
		if f.OptInName(f) == z.Key {
			feature = f
		}
	}
	if feature == nil {
		ui.Error(z.ErrorInvalidKey.With("Key", z.Key))
		return ErrorInvalidKey
	}

	ui.Text(feature.OptInDescription(feature))
	feature.OptInCommit(false)
	if err := c.Feature().OptInUpdate(feature); err != nil {
		ui.Error(z.ErrorUnableToDisableFeature.With("Key", z.Key))
		return err
	}
	ui.Info(z.InfoOptOut.With("Key", z.Key))
	return nil
}

func (z *Disable) Test(c app_control.Control) error {
	if err := rc_exec.Exec(c, &Disable{}, func(r rc_recipe.Recipe) {
		f := &SampleFeatureNotInCatalogue{}
		m := r.(*Disable)
		m.Key = f.OptInName(f)
	}); err != ErrorInvalidKey {
		return ErrorInvalidKey
	}

	if err := rc_exec.Exec(c, &Disable{}, func(r rc_recipe.Recipe) {
		f := &SampleFeature{}
		m := r.(*Disable)
		m.Key = f.OptInName(f)
	}); err != nil {
		return err
	}
	return nil
}