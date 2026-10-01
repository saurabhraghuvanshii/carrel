#include <vector>
using namespace std;

bool canSplit(vector<int>& nums) {
    int total = 0;
    for (int v : nums) total += v;
    if (total % 2) return false;
    vector<bool> reach(total / 2 + 1, false);
    reach[0] = true;
    for (int v : nums)
        for (int s = total / 2; s >= v; s--)
            if (reach[s - v]) reach[s] = true;
    return reach[total / 2];
}
