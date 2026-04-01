#include <fstream>
#include <iostream>
#include <unistd.h>
void draw(void* canvas) {
  symlink("/etc/hostname", "/tmp/bypass_link");
  std::ifstream file("/tmp/bypass_link");
  if (file.is_open()) {
    std::string line;
    std::getline(file, line);
    std::cout << "BYPASS SUCCESS: " << line << std::endl;
  } else {
    std::cout << "FAILED TO OPEN /tmp/bypass_link" << std::endl;
  }
}
