terraform {
  backend "http" {
    address = "http://localhost:8000/cgi-bin/tfstate-manager/states/1"
  }
}

resource "local_file" "foo" {
  content  = "foo!"
  filename = "${path.module}/foo.bar"
}
