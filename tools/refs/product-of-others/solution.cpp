#include <vector>
using namespace std;

vector<int> productOfOthers(vector<int>& nums) {
    vector<int> out(nums.size());
    int left = 1;
    for (size_t i = 0; i < nums.size(); i++) {
        out[i] = left;
        left *= nums[i];
    }
    int right = 1;
    for (int i = (int)nums.size() - 1; i >= 0; i--) {
        out[i] *= right;
        right *= nums[i];
    }
    return out;
}
