#include <string>
using namespace std;

long long squaredDistance(string& steps) {
    long long x = 0, y = 0;
    for (char s : steps) {
        if (s == 'N') y++;
        else if (s == 'S') y--;
        else if (s == 'E') x++;
        else x--;
    }
    return x * x + y * y;
}
