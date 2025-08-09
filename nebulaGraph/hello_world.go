package nebulagraph

import (
	nebula "github.com/vesoft-inc/nebula-go/v3"
)

type Person struct {
	Name     string  `nebula:"name"`
	Age      int     `nebula:"age"`
	Likeness float64 `nebula:"likeness"`
}

func FirstRun() {
	hostAddress := nebula.HostAddress{Host: "127.0.0.1", Port: 3699}

	config, err := nebula.NewSessionPoolConf(
		"root",
		"nebula",
		[]nebula.HostAddress{hostAddress},
		"test",
	)
	if err != nil {
		panic(err)
	}

	sessionPool, err := nebula.NewSessionPool(*config, nebula.DefaultLogger{})
	if err != nil {
		panic(err)
	}

	query := `GO FROM 'Bob' OVER like YIELD
      $^.person.name AS name,
      $^.person.age AS age,
      like.likeness AS likeness`

	resultSet, err := sessionPool.Execute(query)
	if err != nil {
		panic(err)
	}

	var personList []Person
	resultSet.Scan(&personList)
}
