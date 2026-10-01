#include <algorithm>
#include <vector>
using namespace std;

vector<int> firstAndLast(vector<int>& nums, int target) {
    auto lo = lower_bound(nums.begin(), nums.end(), target);
    if (lo == nums.end() || *lo != target) return {-1, -1};
    auto hi = upper_bound(nums.begin(), nums.end(), target);
    return {(int)(lo - nums.begin()), (int)(hi - nums.begin()) - 1};
}
