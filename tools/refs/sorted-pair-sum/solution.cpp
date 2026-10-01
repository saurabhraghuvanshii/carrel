#include <vector>
using namespace std;

vector<int> sortedPair(vector<int>& nums, int target) {
    int i = 0, j = (int)nums.size() - 1;
    while (i < j) {
        int sum = nums[i] + nums[j];
        if (sum == target) return {i, j};
        if (sum < target) i++;
        else j--;
    }
    return {};
}
