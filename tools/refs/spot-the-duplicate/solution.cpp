#include <unordered_set>
#include <vector>
using namespace std;

bool hasDuplicate(vector<int>& nums) {
    unordered_set<int> seen;
    for (int v : nums) {
        if (!seen.insert(v).second) return true;
    }
    return false;
}
