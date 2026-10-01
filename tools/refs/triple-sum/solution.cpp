#include <algorithm>
#include <vector>
using namespace std;

vector<vector<int>> tripleSum(vector<int>& nums, int target) {
    vector<int> s = nums;
    sort(s.begin(), s.end());
    vector<vector<int>> out;
    for (int i = 0; i + 2 < (int)s.size(); i++) {
        if (i > 0 && s[i] == s[i - 1]) continue;
        int lo = i + 1, hi = (int)s.size() - 1;
        while (lo < hi) {
            int sum = s[i] + s[lo] + s[hi];
            if (sum < target) lo++;
            else if (sum > target) hi--;
            else {
                out.push_back({s[i], s[lo], s[hi]});
                while (lo < hi && s[lo] == s[lo + 1]) lo++;
                lo++;
                hi--;
            }
        }
    }
    return out;
}
