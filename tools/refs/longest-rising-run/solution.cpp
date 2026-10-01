#include <algorithm>
#include <vector>
using namespace std;

int longestRising(vector<int>& nums) {
    vector<int> tails;
    for (int v : nums) {
        auto it = lower_bound(tails.begin(), tails.end(), v);
        if (it == tails.end()) tails.push_back(v);
        else *it = v;
    }
    return tails.size();
}
