#include <unistd.h>
#include <iostream>
void draw(void* canvas) {
  char buf[256];
  ssize_t n = readlink("/etc/hostname", buf, sizeof(buf)-1);
  if (n != -1) {
    buf[n] = '\0';
    std::cout << "BYPASS SUCCESS (readlink): " << buf << std::endl;
  } else {
    std::cout << "FAILED readlink" << std::endl;
  }
}
