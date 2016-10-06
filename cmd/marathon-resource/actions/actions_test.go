package actions

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ckaznocha/marathon-resource/cmd/marathon-resource/marathon"
	"github.com/ckaznocha/marathon-resource/cmd/marathon-resource/mocks"
	gomarathon "github.com/gambol99/go-marathon"
	"github.com/golang/mock/gomock"
)

func TestOut(t *testing.T) {
	var (
		ctrl           = gomock.NewController(t)
		mockMarathoner = mocks.NewMockMarathoner(ctrl)
	)
	defer ctrl.Finish()

	gomock.InOrder(
		mockMarathoner.EXPECT().UpdateApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{DeploymentID: "foo", Version: "bar"}, nil),
		mockMarathoner.EXPECT().UpdateApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{}, errors.New("Something went wrong")),
		mockMarathoner.EXPECT().UpdateApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{DeploymentID: "baz", Version: "bar"}, nil),
		mockMarathoner.EXPECT().UpdateApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{DeploymentID: "quux", Version: "bar"}, nil),
		mockMarathoner.EXPECT().UpdateApp(gomock.Any()).Times(1).Return(gomarathon.DeploymentID{DeploymentID: "zork", Version: "bar"}, nil),
	)
	gomock.InOrder(
		mockMarathoner.EXPECT().CheckDeployment("foo").Times(1).Return(false, nil),
		mockMarathoner.EXPECT().CheckDeployment("baz").Times(2).Return(true, nil),
		mockMarathoner.EXPECT().CheckDeployment("quux").Times(1).Return(false, errors.New("Something bad happend")),
		mockMarathoner.EXPECT().CheckDeployment("zork").Times(2).Return(true, nil),
	)
	gomock.InOrder(
		mockMarathoner.EXPECT().DeleteDeployment("baz").Times(1).Return(nil),
		mockMarathoner.EXPECT().DeleteDeployment("zork").Times(1).Return(errors.New("Now way!")),
	)

	type args struct {
		input       InputJSON
		appJSONPath string
		apiclient   marathon.Marathoner
	}
	tests := []struct {
		name    string
		args    args
		want    IOOutput
		wantErr bool
	}{
		{
			"Works",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{Version: Version{Ref: "bar"}},
			false,
		},
		{
			"Bad app json file",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "ajson", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Bad app json file",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app_bad.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Error from UpdateApp",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Deployment times out",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Check deployment errors",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
		{
			"Delete deployment errors",
			args{
				input: InputJSON{
					Params: Params{AppJSON: "app.json", TimeOut: 2},
					Source: Source{},
				},
				appJSONPath: "../fixtures",
				apiclient:   mockMarathoner,
			},
			IOOutput{},
			true,
		},
	}

	for _, tt := range tests {
		got, err := Out(tt.args.input, tt.args.appJSONPath, tt.args.apiclient)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Out() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Out() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
