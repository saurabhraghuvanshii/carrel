#include <algorithm>
#include <vector>
using namespace std;

int worstBurst(vector<int>& times, int window) {
    int best = 0, lo = 0;
    for (int hi = 0; hi < (int)times.size(); hi++) {
        while (times[lo] <= times[hi] - window) lo++;
        best = max(best, hi - lo + 1);
    }
    return best;
}
