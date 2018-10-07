resource "teamcity_project" "parent" {
    name = "Parent"
    description = "Parent Project, will be created under the 'Root' project"
}

resource "teamcity_project" "child" {
    name = "Child"
    description = "Child Project, will be created under 'Parent' project"
    parent_id = "${teamcity_project.parent.id}"

    config_params = {
        variable1 = "config_value1"
        variable2 = "config_value2"
    }

    env_params = {
        variable1 = "env_value1"
        variable2 = "env_value2"
    }

    sys_params = { 
        variable1 = "system_value1"
        variable2 = "system_value2"
    }
}