#include <deque>
#include <vector>
using namespace std;

vector<int> windowMaximums(vector<int>& nums, int k) {
    vector<int> out;
    deque<int> queue;
    for (int i = 0; i < (int)nums.size(); i++) {
        while (!queue.empty() && nums[queue.back()] <= nums[i]) queue.pop_back();
        queue.push_back(i);
        if (queue.front() <= i - k) queue.pop_front();
        if (i >= k - 1) out.push_back(nums[queue.front()]);
    }
    return out;
}
