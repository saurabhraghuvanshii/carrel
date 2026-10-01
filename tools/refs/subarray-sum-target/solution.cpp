#include <unordered_map>
#include <vector>
using namespace std;

int countSubarrays(vector<int>& nums, int target) {
    unordered_map<int, int> seen{{0, 1}};
    int sum = 0, count = 0;
    for (int v : nums) {
        sum += v;
        auto it = seen.find(sum - target);
        if (it != seen.end()) count += it->second;
        seen[sum]++;
    }
    return count;
}
