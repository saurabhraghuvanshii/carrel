#include <algorithm>
#include <vector>
using namespace std;

int shortestWalk(vector<int>& stops, int start) {
    int lo = start, hi = start;
    for (int s : stops) lo = min(lo, s), hi = max(hi, s);
    int left = start - lo, right = hi - start;
    return left + right + min(left, right);
}
