#include <algorithm>
#include <unordered_map>
#include <vector>
using namespace std;

int longestTwoKinds(vector<int>& kinds) {
    unordered_map<int, int> count;
    int best = 0, start = 0;
    for (int i = 0; i < (int)kinds.size(); i++) {
        count[kinds[i]]++;
        while (count.size() > 2) {
            if (--count[kinds[start]] == 0) count.erase(kinds[start]);
            start++;
        }
        best = max(best, i - start + 1);
    }
    return best;
}
