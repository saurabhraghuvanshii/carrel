#include <algorithm>
#include <vector>
using namespace std;

int smallestCapacity(vector<int>& weights, int k) {
    int lo = 0, hi = 0;
    for (int w : weights) lo = max(lo, w), hi += w;
    while (lo < hi) {
        int mid = lo + (hi - lo) / 2, count = 1, load = 0;
        for (int w : weights) {
            if (load + w > mid) count++, load = 0;
            load += w;
        }
        if (count <= k) hi = mid;
        else lo = mid + 1;
    }
    return lo;
}
