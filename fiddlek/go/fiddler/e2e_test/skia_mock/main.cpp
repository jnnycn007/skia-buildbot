#include <iostream>
void draw(void*);
int main() {
    draw(nullptr);
    std::cout << "{\"Raster\": \"\", \"Gpu\": \"\"}" << std::endl;
    return 0;
}
