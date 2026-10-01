#include <string>
#include <vector>
using namespace std;

string commonStart(vector<string>& words) {
    size_t k = 0;
    while (true) {
        for (auto& w : words)
            if (k >= w.size() || w[k] != words[0][k]) return words[0].substr(0, k);
        k++;
    }
}
