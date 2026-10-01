#include <algorithm>
#include <functional>
#include <vector>
using namespace std;

long long cutCost(int rows, int cols, vector<int>& horizontal, vector<int>& vertical) {
    sort(horizontal.rbegin(), horizontal.rend());
    sort(vertical.rbegin(), vertical.rend());
    size_t h = 0, v = 0;
    long long total = 0;
    while (h < horizontal.size() || v < vertical.size()) {
        if (v == vertical.size() || (h < horizontal.size() && horizontal[h] >= vertical[v])) total += (long long)horizontal[h++] * (v + 1);
        else total += (long long)vertical[v++] * (h + 1);
    }
    return total;
}
