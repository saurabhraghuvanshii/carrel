#include <vector>
using namespace std;

int loneNumber(vector<int>& nums) {
    int x = 0;
    for (int v : nums) x ^= v;
    return x;
}
