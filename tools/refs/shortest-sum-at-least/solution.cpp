#include <vector>
using namespace std;

int shortestAtLeast(vector<int>& nums, int target) {
    int best = 0, sum = 0, start = 0;
    for (int i = 0; i < (int)nums.size(); i++) {
        sum += nums[i];
        while (sum >= target) {
            if (best == 0 || i - start + 1 < best) best = i - start + 1;
            sum -= nums[start++];
        }
    }
    return best;
}
