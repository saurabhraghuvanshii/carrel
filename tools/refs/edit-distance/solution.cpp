#include <algorithm>
#include <string>
#include <vector>
using namespace std;

int fewestEdits(string& first, string& second) {
    size_t m = second.size();
    vector<int> prev(m + 1), cur(m + 1);
    for (size_t j = 0; j <= m; j++) prev[j] = j;
    for (size_t i = 1; i <= first.size(); i++) {
        cur[0] = i;
        for (size_t j = 1; j <= m; j++)
            cur[j] = first[i - 1] == second[j - 1] ? prev[j - 1] : 1 + min({prev[j - 1], prev[j], cur[j - 1]});
        swap(prev, cur);
    }
    return prev[m];
}
