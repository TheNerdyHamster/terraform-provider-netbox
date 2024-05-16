package netbox

import "fmt"

func testAccNetboxVirtualDiskDataSourceDependencies(testName string) string {
    	return fmt.Sprintf(`
resource "netbox_tag" "tag_a" {
        name = "[%[1]s_a]"
        color_hex = "123456"
}

resource "netbox_site" "test" {
        name = "%[1]s"
        status = active
}

resource "netbox_virtual_machine" "test0" {
        name = "%[1]s_0"
        site_id = netbox_site.test.id
}

resource "netbox_virtual_machine" "test1" {
        name = "%[1]s_1"
        site_id = netbox_site.test.id
}

resource "netbox_virtual_disk" "test0" {
        name = "%[1]s"
        description = "description"
        size_gb = 30
        virtual_machine_id = netbox_virtual_machine.test0.id
        tags = [netbox_tag.tag_a.name]
}

resource "netbox_virtual_disk" "test1" {
        name = "%[1]s_0"
        description = "description"
        size_gb = 30
        virtual_machine_id = netbox_virtual_machine.test0.id
        tags = [netbox_tag.tag_a.name]
}

resource "netbox_virtual_disk" "test2" {
        name = "%[1]s"
        description = "description1"
        size_gb = 30
        virtual_machine_id = netbox_virtual_machine.test1.id
}`, testName)
}

const testAccNetboxVirtualDiskDataSourceFilterVM = `
data "netbox_virtual_disk" "test" {
    filter {
        name = "vm_id"
        value = netbox_virtual_machine.test1.id
    }
}
`
