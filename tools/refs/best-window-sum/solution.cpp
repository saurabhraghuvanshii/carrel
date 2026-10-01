#include <algorithm>
#include <vector>
using namespace std;

int bestWindowSum(vector<int>& nums, int k) {
    int sum = 0;
    for (int i = 0; i < k; i++) sum += nums[i];
    int best = sum;
    for (int i = k; i < (int)nums.size(); i++) {
        sum += nums[i] - nums[i - k];
        best = max(best, sum);
    }
    return best;
}
