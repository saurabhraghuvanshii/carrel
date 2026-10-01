#include <algorithm>
#include <unordered_set>
#include <vector>
using namespace std;

int longestRun(vector<int>& nums) {
    unordered_set<int> set(nums.begin(), nums.end());
    int best = 0;
    for (int v : set) {
        if (set.count(v - 1)) continue;
        int length = 1;
        while (set.count(v + length)) length++;
        best = max(best, length);
    }
    return best;
}
