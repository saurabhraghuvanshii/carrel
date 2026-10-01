#include <algorithm>
#include <string>
#include <vector>
using namespace std;

int longestShared(string& first, string& second) {
    vector<int> prev(second.size() + 1, 0), cur(second.size() + 1, 0);
    for (size_t i = 1; i <= first.size(); i++) {
        for (size_t j = 1; j <= second.size(); j++)
            cur[j] = first[i - 1] == second[j - 1] ? prev[j - 1] + 1 : max(prev[j], cur[j - 1]);
        swap(prev, cur);
    }
    return prev[second.size()];
}
