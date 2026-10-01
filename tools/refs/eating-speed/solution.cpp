#include <algorithm>
#include <vector>
using namespace std;

int slowestSpeed(vector<int>& jobs, int h) {
    int lo = 1, hi = *max_element(jobs.begin(), jobs.end());
    while (lo < hi) {
        int mid = lo + (hi - lo) / 2;
        long long hours = 0;
        for (int j : jobs) hours += (j + (long long)mid - 1) / mid;
        if (hours <= h) hi = mid;
        else lo = mid + 1;
    }
    return lo;
}
