variable "my_ip" {
  description = "My pc public IP address"
  type        = string
  sensitive   = true
}


variable "gh_runner_token" {
  description = "The 60-minute registration token from GitHub for selfhosted runnner"
  type        = string
  sensitive   = true # hide from terminal logs
}
