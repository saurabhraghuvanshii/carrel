#include <vector>
using namespace std;

int removeRepeats(vector<int>& nums) {
    int k = 0;
    for (size_t i = 0; i < nums.size(); i++) {
        if (k == 0 || nums[k - 1] != nums[i]) nums[k++] = nums[i];
    }
    return k;
}
