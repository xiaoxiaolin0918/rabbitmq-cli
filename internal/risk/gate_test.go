package risk

import (
	"reflect"
	"testing"
)

func TestHighRisk(t *testing.T) {
	cases := []struct {
		args []string
		want bool
	}{
		{[]string{"show", "overview"}, false},
		{[]string{"show", "churn"}, false},
		{[]string{"queues", "list"}, false},
		{[]string{"list", "queues"}, false},
		{[]string{"list", "consumers"}, false},
		{[]string{"users", "list"}, false},
		{[]string{"users", "permissions"}, false},
		{[]string{"users", "connections"}, false},
		{[]string{"policies", "list_matching_object"}, false},
		{[]string{"health_check", "local_alarms"}, false},
		{[]string{"connections", "list_of_user", "--username", "x"}, false},
		{[]string{"--non-interactive", "queues", "list"}, false},
		{[]string{"-N", "lab", "show", "overview"}, false},
		{[]string{"queues", "list", "--name", "delete"}, false},
		{[]string{"config_file", "show_path"}, false},
		{[]string{}, false},
		{[]string{"delete", "queue", "--name", "x"}, true},
		{[]string{"purge", "queue", "--name", "x"}, true},
		{[]string{"queues", "delete", "--name", "x"}, true},
		{[]string{"definitions", "import", "f.json"}, true},
		{[]string{"definitions", "export"}, true},
		{[]string{"definitions", "import_into_vhost", "--file", "f.json"}, true},
		{[]string{"export", "--file", "-"}, true},
		{[]string{"get", "--queue", "q"}, true},
		// leading / interleaved global flags must not bypass the gate
		{[]string{"--vhost", "/", "declare", "queue", "--name", "x"}, true},
		{[]string{"-V", "/", "delete", "queue", "--name", "x"}, true},
		{[]string{"--non-interactive", "delete", "queue", "--name", "x"}, true},
		{[]string{"--vhost=/", "purge", "queue", "--name", "x"}, true},
		{[]string{"queues", "-V", "/", "delete", "--name", "x"}, true},
		{[]string{"-q", "vhosts", "delete", "--name", "x"}, true},
		// group verbs previously missed
		{[]string{"queues", "delete_multiple", "--name-pattern", "x"}, true},
		{[]string{"vhosts", "delete_multiple", "--name-pattern", "x"}, true},
		{[]string{"vhosts", "enable_deletion_protection", "--name", "x"}, true},
		{[]string{"exchanges", "bind", "--source", "a"}, true},
		{[]string{"exchanges", "unbind", "--source", "a"}, true},
		{[]string{"policies", "set", "--name", "p"}, true},
		{[]string{"policies", "patch", "--name", "p"}, true},
		{[]string{"policies", "declare_override", "--name", "p"}, true},
		{[]string{"policies", "update_definitions_of_all", "--definition-key", "k"}, true},
		{[]string{"operator_policies", "declare", "--name", "p"}, true},
		{[]string{"parameters", "set", "--name", "p"}, true},
		{[]string{"global_parameters", "clear", "--name", "p"}, true},
		{[]string{"feature_flags", "enable", "--name", "f"}, true},
		{[]string{"feature_flags", "enable_all"}, true},
		{[]string{"federation", "declare_upstream", "--name", "u"}, true},
		{[]string{"shovels", "declare_amqp091", "--name", "s"}, true},
		{[]string{"connections", "close_of_user", "--username", "u"}, true},
		{[]string{"user_limits", "declare", "--user", "u"}, true},
		{[]string{"vhost_limits", "delete", "--name", "max-queues"}, true},
		{[]string{"config_file", "delete_node", "--node", "n"}, true},
		{[]string{"config_file", "add_node", "--node", "n"}, true},
		{[]string{"rebalance", "queues"}, true},
		{[]string{"queues", "rebalance"}, true},
		{[]string{"tanzu", "wsr", "enable"}, true},
	}
	for _, c := range cases {
		if got := HighRisk(c.args); got != c.want {
			t.Errorf("%v: got %v want %v", c.args, got, c.want)
		}
	}
}

func TestRedact(t *testing.T) {
	in := []string{"declare", "user", "--name", "u", "--password", "s3cret", "-p", "x", "--password=y"}
	got := Redact(in)
	want := []string{"declare", "user", "--name", "u", "--password", "***", "-p", "***", "--password=***"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
	if in[5] != "s3cret" {
		t.Fatal("Redact must not mutate input")
	}
}
