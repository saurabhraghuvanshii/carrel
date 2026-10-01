#include <functional>
#include <queue>
#include <vector>
using namespace std;

int kthLargest(vector<int>& nums, int k) {
    priority_queue<int, vector<int>, greater<int>> top;
    for (int v : nums) {
        top.push(v);
        if ((int)top.size() > k) top.pop();
    }
    return top.top();
}
