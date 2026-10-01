#include <algorithm>
#include <vector>
using namespace std;

int maxCollect(vector<int>& nums) {
    int take = 0, skip = 0;
    for (int v : nums) {
        int next = skip + v;
        skip = max(take, skip);
        take = next;
    }
    return max(take, skip);
}
