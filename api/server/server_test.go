package server

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	"github.com/quilt/quilt/api/pb"
	"github.com/quilt/quilt/db"
	"github.com/quilt/quilt/stitch"
	"github.com/stretchr/testify/assert"
)

func checkQuery(t *testing.T, s server, table db.TableType, exp string) {
	reply, err := s.Query(context.Background(),
		&pb.DBQuery{Table: string(table)})

	assert.NoError(t, err)
	assert.Equal(t, exp, reply.TableContents, "Wrong query response")
}

func TestMachineResponse(t *testing.T) {
	t.Parallel()

	conn := db.New()
	conn.Txn(db.AllTables...).Run(func(view db.Database) error {
		m := view.InsertMachine()
		m.Role = db.Master
		m.Provider = db.Amazon
		m.Size = "size"
		m.PublicIP = "8.8.8.8"
		m.PrivateIP = "9.9.9.9"
		view.Commit(m)

		return nil
	})

	exp := `[{"ID":1,"StitchID":"","Role":"Master","Provider":"Amazon","Region":"",` +
		`"Size":"size","DiskSize":0,"SSHKeys":null,"FloatingIP":"",` +
		`"Preemptible":false,"CloudID":"","PublicIP":"8.8.8.8",` +
		`"PrivateIP":"9.9.9.9","Connected":false}]`

	checkQuery(t, server{conn}, db.MachineTable, exp)
}

func TestContainerResponse(t *testing.T) {
	t.Parallel()

	conn := db.New()
	conn.Txn(db.AllTables...).Run(func(view db.Database) error {
		c := view.InsertContainer()
		c.DockerID = "docker-id"
		c.Image = "image"
		c.Command = []string{"cmd", "arg"}
		c.Labels = []string{"labelA", "labelB"}
		view.Commit(c)

		return nil
	})

	exp := `[{"DockerID":"docker-id","Command":["cmd","arg"],` +
		`"Labels":["labelA","labelB"],"Created":"0001-01-01T00:00:00Z",` +
		`"Image":"image"}]`

	checkQuery(t, server{conn}, db.ContainerTable, exp)
}

func TestBadDeployment(t *testing.T) {
	conn := db.New()
	s := server{conn: conn}

	badDeployment := `{`

	_, err := s.Deploy(context.Background(),
		&pb.DeployRequest{Deployment: badDeployment})

	assert.EqualError(t, err, "unexpected end of JSON input")
}
func TestInvalidImage(t *testing.T) {
	conn := db.New()
	s := server{conn: conn}
	testInvalidImage(t, s, "has:morethan:two:colons",
		"could not parse container image has:morethan:two:colons: "+
			"invalid reference format")
	testInvalidImage(t, s, "has-empty-tag:",
		"could not parse container image has-empty-tag:: "+
			"invalid reference format")
	testInvalidImage(t, s, "has-empty-tag::digest",
		"could not parse container image has-empty-tag::digest: "+
			"invalid reference format")
	testInvalidImage(t, s, "hasCapital",
		"could not parse container image hasCapital: "+
			"invalid reference format: repository name must be lowercase")
}

func testInvalidImage(t *testing.T, s server, img, expErr string) {
	deployment := fmt.Sprintf(`
	{"Containers":[
		{"ID": "1",
                "Image": {"Name": "%s"},
                "Command":[
                        "sleep",
                        "10000"
                ],
                "Env": {}
	}]}`, img)

	_, err := s.Deploy(context.Background(),
		&pb.DeployRequest{Deployment: deployment})
	assert.EqualError(t, err, expErr)
}

func TestDeploy(t *testing.T) {
	conn := db.New()
	s := server{conn: conn}

	createMachineDeployment := `
	{"Machines":[
		{"Provider":"Amazon",
		"Role":"Master",
		"Size":"m4.large"
	}, {"Provider":"Amazon",
		"Role":"Worker",
		"Size":"m4.large"
	}]}`

	_, err := s.Deploy(context.Background(),
		&pb.DeployRequest{Deployment: createMachineDeployment})

	assert.NoError(t, err)

	var blueprint string
	conn.Txn(db.AllTables...).Run(func(view db.Database) error {
		clst, err := view.GetCluster()
		assert.NoError(t, err)
		blueprint = clst.Blueprint
		return nil
	})

	exp, err := stitch.FromJSON(createMachineDeployment)
	assert.NoError(t, err)

	actual, err := stitch.FromJSON(blueprint)
	assert.NoError(t, err)

	assert.Equal(t, exp, actual)
}

func TestVagrantDeployment(t *testing.T) {
	conn := db.New()
	s := server{conn: conn}

	vagrantDeployment := `
	{"Machines":[
		{"Provider":"Vagrant",
		"Role":"Master",
		"Size":"m4.large"
	}, {"Provider":"Vagrant",
		"Role":"Worker",
		"Size":"m4.large"
	}]}`
	vagrantErrMsg := "The Vagrant provider is in development." +
		" The stitch will continue to run, but" +
		" probably won't work correctly."

	_, err := s.Deploy(context.Background(),
		&pb.DeployRequest{Deployment: vagrantDeployment})

	assert.Error(t, err, vagrantErrMsg)

	var blueprint string
	conn.Txn(db.AllTables...).Run(func(view db.Database) error {
		clst, err := view.GetCluster()
		assert.NoError(t, err)
		blueprint = clst.Blueprint
		return nil
	})

	exp, err := stitch.FromJSON(vagrantDeployment)
	assert.NoError(t, err)

	actual, err := stitch.FromJSON(blueprint)
	assert.NoError(t, err)

	assert.Equal(t, exp, actual)
}
